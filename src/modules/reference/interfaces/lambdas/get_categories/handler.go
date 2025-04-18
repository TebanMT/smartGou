package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TebanMT/smartGou/infraestructure/db"
	app "github.com/TebanMT/smartGou/src/modules/reference/app/categories"
	"github.com/TebanMT/smartGou/src/modules/reference/infrastructure/repositories"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	"github.com/TebanMT/smartGou/src/shared/criteria"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/middleware"
	"github.com/TebanMT/smartGou/src/shared/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type Request struct {
	MetaCategoryID *uuid.UUID `schema:"meta_category_id"`
	NameLike       *string    `schema:"name_like"`
	Limit          *int       `schema:"limit"`
	Offset         *int       `schema:"offset"`
	Paged          *bool      `schema:"paged"`
	OrderBy        *string    `schema:"order_by"`
	OrderDir       *string    `schema:"order_dir"`
}

type Category struct {
	ID             string `json:"id"`
	NameEn         string `json:"name_en"`
	NameEs         string `json:"name_es"`
	Icon           string `json:"icon"`
	Color          string `json:"color"`
	Description    string `json:"description"`
	MetaCategoryID string `json:"meta_category_id"`
}

type ResponsePaginated struct {
	Total  int        `json:"total"`
	Offset *int       `json:"offset"`
	Limit  *int       `json:"limit"`
	Data   []Category `json:"categories"`
}

var cognitoService *cognito.CognitoService
var err error
var getCategoriesUseCase *app.GetCategoriesUseCase

func init() {
	dbInstance := db.InitConnection()
	unitOfWork := sharedDomain.NewUnitOfWork(dbInstance)
	categoriesRepository := repositories.NewGormCategoriesRepository(dbInstance)
	getCategoriesUseCase = app.NewGetCategoriesUseCase(categoriesRepository, unitOfWork)
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Get Categories
// @Description  This endpoint gets the categories for a user. It can be filtered by meta category id or name like.
// @Description  If no filters are provided, all categories will be returned.
// @Description  If meta category id is provided, only categories related to that meta category will be returned.
// @Description  If name like is provided, only categories related to that name will be returned.
// @Description  If paged is true, you must provide valid limit and offset values.
// @Description  If paged is false, all categories will be returned, ignoring limit and offset values.
// @Description  If order by is provided, the order of the categories will be based on the order by field and order dir field will be required.
// @Description  If order by is not provided, the categories will be returned in the order of the database.
// @Decription You can provide a name like filter to get categories related to that name, the name like filter will be applied only on name_en (english name).
// @Tags reference
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param meta_category_id query string false "Meta Category ID"
// @Param name_like query string false "Name Like"
// @Param order_by query string false "Order By"
// @Param order_dir query string false "Order Dir"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param paged query bool false "Paged"
// @Success 200 {object} utils.Response[ResponsePaginated]
// @Success 206 {object} utils.Response[[]Category] "The code 206 just means the response has not been paginated but the response will send with 200 code"
// @Failure 400 {object} utils.Response[any]
// @Failure 401 {object} utils.Response[any]
// @Failure 404 {object} utils.Response[any]
// @Failure 500 {object} utils.Response[any]
// @Router /references/categories [get]
func getCategoriesLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var params Request
	var response ResponsePaginated
	err := utils.ValidateQueryStringParameters(request, &params)
	if err != nil {
		fmt.Println("error", err)
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}
	fmt.Println("params", params)

	criteria := &app.CategoryQuery{
		BaseCriteria: criteria.BaseCriteria{
			Limit:    params.Limit,
			Offset:   params.Offset,
			OrderBy:  params.OrderBy,
			OrderDir: params.OrderDir,
			Paged:    params.Paged,
		},
		MetaCategoryID: params.MetaCategoryID,
		NameLike:       params.NameLike,
	}
	fmt.Println(criteria.Debug())
	categories, err := getCategoriesUseCase.GetCategories(ctx, criteria)
	httpCode := utils.DomainErrorToHttpCode(err)
	switch true {
	case err == nil:
		data := assembleCategoriesResponse(categories)
		if params.Paged != nil && *params.Paged {
			response = ResponsePaginated{
				Total:  *criteria.TotalOfRecords,
				Offset: params.Offset,
				Limit:  params.Limit,
				Data:   data,
			}
			return utils.JsonResponse(200, "", response, "")
		}
		return utils.JsonResponse(200, "", data, "")
	case httpCode != nil && *httpCode != 0:
		return utils.JsonResponse[any](*httpCode, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func assembleCategoriesResponse(categories []sharedDomain.Category) []Category {
	data := []Category{}
	for _, category := range categories {
		data = append(data, Category{
			ID:             category.ID.String(),
			NameEn:         category.NameEn,
			NameEs:         category.NameEs,
			Icon:           category.Icon,
			Color:          category.Color,
			Description:    category.Description,
			MetaCategoryID: category.MetaCategoryID,
		})
	}
	return data
}

func main() {
	handler := middleware.AuthenticationMiddleware(cognitoService, getCategoriesLambdaHandler)
	lambda.Start(handler)
}
