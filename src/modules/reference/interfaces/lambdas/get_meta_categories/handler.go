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

type Parameters struct {
	NameLike *string `schema:"name_like"`
	Limit    *int    `schema:"limit"`
	Offset   *int    `schema:"offset"`
	Paged    *bool   `schema:"paged"`
	OrderDir *string `schema:"order_dir"`
	OrderBy  *string `schema:"order_by"`
}

type MetaCategoryResponse struct {
	ID          uuid.UUID
	NameEn      string
	NameEs      string
	Icon        string
	Color       string
	Description string
}

type ResponsePaginated struct {
	Total  int                    `json:"total"`
	Offset *int                   `json:"offset"`
	Limit  *int                   `json:"limit"`
	Data   []MetaCategoryResponse `json:"meta_categories"`
}

var cognitoService *cognito.CognitoService
var err error
var getMetaCategoriesUseCase *app.GetMetaCategoriesUseCase

func init() {
	dbInstance := db.InitConnection()
	unitOfWork := sharedDomain.NewUnitOfWork(dbInstance)
	categoriesRepository := repositories.NewGormCategoriesRepository(dbInstance)
	getMetaCategoriesUseCase = app.NewGetMetaCategoriesUseCase(categoriesRepository, unitOfWork)
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Get meta categories
// @Description  This endpoint gets the meta categories catalog.
// @Description  If no filters are provided, all meta categories will be returned.
// @Description  If name like is provided, only meta categories related to that name will be returned.
// @Description  If paged is true, you must provide valid limit and offset values.
// @Description  If paged is false, all meta categories will be returned, ignoring limit and offset values.
// @Description  If order by is provided, the order of the meta categories will be based on the order by field and order dir field will be required.
// @Description  If order by is not provided, the meta categories will be returned in the order of the database.
// @Decription You can provide a name like filter to get meta categories related to that name, the name like filter will be applied only on name_en (english name).
// @Tags reference
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param name_like query string false "Name Like"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param paged query bool false "Paged"
// @Param order_by query string false "Order By"
// @Param order_dir query string false "Order Dir"
// @Success 200 {object} utils.Response[ResponsePaginated]
// @Success 206 {object} utils.Response[[]MetaCategoryResponse] "The code 206 just means the response has not been paginated but the response will send with 200 code"
// @Failure 400 {object} utils.Response[any]
// @Failure 401 {object} utils.Response[any]
// @Failure 500 {object} utils.Response[any]
// @Router /references/meta-categories [get]
func getMetaCategoriesLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var params Parameters
	var response ResponsePaginated
	err := utils.ValidateQueryStringParameters(request, &params)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	criteria := &app.CategoryQuery{
		BaseCriteria: criteria.BaseCriteria{
			Limit:    params.Limit,
			Offset:   params.Offset,
			OrderBy:  params.OrderBy,
			OrderDir: params.OrderDir,
			Paged:    params.Paged,
		},
		NameLike: params.NameLike,
	}
	fmt.Println(criteria.Debug())
	metaCategories, err := getMetaCategoriesUseCase.GetMetaCategories(ctx, criteria)
	httpCode := utils.DomainErrorToHttpCode(err)
	switch true {
	case err == nil:
		data := assembleMetaCategoriesResponse(metaCategories)
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

func assembleMetaCategoriesResponse(metaCategories []sharedDomain.MetaCategory) []MetaCategoryResponse {
	data := []MetaCategoryResponse{}
	for _, metaCategory := range metaCategories {
		data = append(data, MetaCategoryResponse{
			ID:          metaCategory.ID,
			NameEn:      metaCategory.NameEn,
			NameEs:      metaCategory.NameEs,
			Icon:        metaCategory.Icon,
			Color:       metaCategory.Color,
			Description: metaCategory.Description,
		})
	}
	return data
}

func main() {
	handler := middleware.AuthenticationMiddleware(cognitoService, getMetaCategoriesLambdaHandler)
	lambda.Start(handler)
}
