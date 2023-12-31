package controller

import (
	"account/api"
	queries "account/config/sqlc"
	"account/internal/database"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nav-inc/datetime"

	"github.com/gin-gonic/gin"
)

type ItemController struct{}

func (ctrl *ItemController) Get(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// CreateItem
//
//	@Summary	创建账目
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Param		amount		body		int				true	"金额（单位：分）"	example(100)
//	@Param		kind		body		queries.Kind	true	"类型"		example(expenses)
//	@Param		happened_at	body		string			true	"发生时间"		example(2023-09-26T00:00:00Z)
//	@Param		tag_ids		body		[]int			true	"标签ID列表"	example([1,2,3])
//
//	@Success	200			{object}	api.CreateItemResponse
//	@Failure	401			{string}	string	无效的JWT
//	@Failure	422			{string}	string	参数错误
//	@Router		/api/v1/items [post]
func (ctrl *ItemController) Create(c *gin.Context) {
	var body api.CreateItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(422, "参数错误", err)
		return
	}
	// 获取当前登录用户
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	item, err := q.CreateItem(c, queries.CreateItemParams{
		UserID:     user.ID,
		Amount:     body.Amount,
		Kind:       body.Kind,
		HappenedAt: body.HappenedAt,
		TagIds:     body.TagIds,
	})
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": item,
	})

}

func (ctrl *ItemController) Update(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) GetBanlance(c *gin.Context) {
	query := c.Request.URL.Query()
	happenedAfterStr := query.Get("happened_after")
	happenedBeforeStr := query.Get("happened_before")
	happenedAfter, err := datetime.Parse(happenedAfterStr, time.Local)
	if err != nil {
		happenedAfter = time.Now().AddDate(-100, 0, 0)
	}
	happenedBefore, err := datetime.Parse(happenedBeforeStr, time.Local)
	if err != nil {
		happenedBefore = time.Now().AddDate(1, 0, 0)
	}

	q := database.NewQuery()
	items, err := q.ListItemsHappenedBetween(c, queries.ListItemsHappenedBetweenParams{
		HappenedAfter:  happenedAfter,
		HappenedBefore: happenedBefore,
	})
	if err != nil {
		log.Printf("list items error: %v", err)
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	// 计算返回值
	var r api.GetBalanceResponse
	for _, item := range items {
		if item.Kind == "in_come" {
			r.Income += int(item.Amount)
		}
		if item.Kind == "expenses" {
			r.Expenses += int(item.Amount)
		}
	}
	r.Balance = r.Income - r.Expenses
	c.JSON(http.StatusOK, r)

}

func (ctrl *ItemController) Find(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (ctrl *ItemController) Destory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// GetPagedItems
//
//	@Summary	分页获取账目
//	@Accept		json
//	@Produce	json
//
//	@Security	Bearer
//
//	@Param		page			query		int		false	"页码"	example(1)
//	@Param		page_size		query		int		false	"每页条数"	example(1)
//	@Param		happened_after	query		string	false	"开始时间"
//	@Param		happened_before	query		string	false	"结束时间"
//
//	@Success	200				{object}	api.GetPagedItemsResponse
//	@Failure	500				{string}	string	服务器繁忙
//	@Router		/api/v1/items [get]
func (ctrl *ItemController) GetPaged(c *gin.Context) {
	var params api.GetPagedItemsRequest
	params.Page = 1
	params.PageSize = 10

	query := c.Request.URL.Query()
	if len(query["page"]) > 0 {
		pageStr := query["page"][0]
		if page, err := strconv.Atoi(pageStr); err == nil {
			params.Page = int32(page)
		}
	}
	if len(query["page_size"]) > 0 {
		pageSizeStr := query["page_size"][0]
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			params.PageSize = int32(pageSize)
		}
	}

	// 时间参数处理
	happenedBefore, has := c.Params.Get("happened_before")
	if has {
		if t, err := time.Parse(time.RFC3339, happenedBefore); err == nil {
			params.HappenedBefore = t
		}
	}
	happenedAfter, has := c.Params.Get("happened_after")
	if has {
		if t, err := time.Parse(time.RFC3339, happenedAfter); err == nil {
			params.HappenedAfter = t
		}
	}

	q := database.NewQuery()
	items, err := q.ListItems(c, queries.ListItemsParams{
		Offset: (params.Page - 1) * params.PageSize,
		Limit:  params.PageSize,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	count, err := q.CountItems(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "服务器繁忙")
		return
	}
	c.JSON(http.StatusOK, api.GetPagedItemsResponse{
		Resources: items,
		Pager: api.Pager{
			Page:     params.Page,
			PageSize: params.PageSize,
			Total:    count,
		},
	})
}

func (ctrl *ItemController) GetSummary(c *gin.Context) {
	var query api.GetSummaryRequest
	if err := c.BindQuery(&query); err != nil {
		er := api.ErrorResponse{Errors: map[string][]string{}}
		switch e := err.(type) {
		case validator.ValidationErrors:
			for _, ve := range e {
				// 错误标签
				tag := ve.Tag()
				// 错误字段
				field := ve.Field()
				if er.Errors[field] == nil {
					er.Errors[field] = []string{}
				}
				// 给该field的数组追加一项
				er.Errors[field] = append(er.Errors[field], tag)
			}
			c.JSON(http.StatusUnprocessableEntity, r)
		default:
			c.Writer.WriteString("参数错误")
		}
		return
	}

	me, _ := c.Get("me")
	user, _ := me.(queries.User)

	q := database.NewQuery()
	items, err := q.ListItemsByHappenedAtAndKind(c, queries.ListItemsByHappenedAtAndKindParams{
		HappenedAfter:  query.HappenedAfter,
		HappenedBefore: query.HappenedBefore,
		Kind:           query.Kind,
		UserID:         user.ID,
	})
	if err != nil {
		log.Printf("list items error: %v", err)
		c.String(http.StatusInternalServerError, "数据库内部错误")
		return
	}
	// 按 happened_at 进行分组
	if query.GroupBy == "happened_at" {
		r := api.GetSummaryByHappenedAtResponse{}
		r.Groups = []api.SummaryGroupByHappenedAt{}
		r.Total = 0
		for _, item := range items {
			k := item.HappenedAt.Format("2006-01-02") // 相当于 yyyy-MM-dd
			r.Total += int(item.Amount)
			found := false
			for index, group := range r.Groups {
				if group.HappenedAt == k {
					found = true
					r.Groups[index].Amount += int(item.Amount)
				}
			}
			if !found {
				r.Groups = append(r.Groups, api.SummaryGroupByHappenedAt{
					HappenedAt: k,
					Amount:     int(item.Amount),
				})
			}
		}
		// 对 group 按 happened_at 从前到后进行排序
		sort.Slice(r.Groups, func(i int, j int) bool {
			return r.Groups[i].HappenedAt < r.Groups[j].HappenedAt
		})

		c.JSON(200, r)
	} else {
		// 按 tag_id 进行分组
		r := api.GetSummaryByTagIDResponse{}
		r.Groups = []api.SummaryGroupByTagID{}
		r.Total = 0
		for _, item := range items {
			// 只取第一个 TagID
			k := item.TagIds[0]
			r.Total += int(item.Amount)
			found := false
			for index, group := range r.Groups {
				if group.TagID == k {
					found = true
					r.Groups[index].Amount += int(item.Amount)
				}
			}
			if !found {
				t, err := q.FindTag(c, queries.FindTagParams{
					ID:     k,
					UserID: user.ID,
				})
				if err != nil {
					log.Printf("find tag error: %v", err)
					c.String(http.StatusInternalServerError, "数据库内部错误")
					return
				}
				r.Groups = append(r.Groups, api.SummaryGroupByTagID{
					TagID:  k,
					Tag:    t,
					Amount: int(item.Amount),
				})
			}
		}
		sort.Slice(r.Groups, func(i int, j int) bool {
			return r.Groups[i].TagID < r.Groups[j].TagID
		})

		c.JSON(200, r)
	}
}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("v1")
	// 注册路由
	v1.POST("/items", ctrl.Create)
	v1.GET("/items", ctrl.GetPaged)
	v1.GET("/items/balance", ctrl.GetBanlance)
	v1.GET("/items/summary", ctrl.GetSummary)
}
