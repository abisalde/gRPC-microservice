package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/abisalde/grpc-microservice/auth/pkg/ent/user"
	"github.com/abisalde/grpc-microservice/gateway/internal/graph/model"
)

type AuthDirective struct {
}

func NewAuthDirective() *AuthDirective {
	return &AuthDirective{}
}

func (a *AuthDirective) Auth(ctx context.Context, obj interface{}, next graphql.Resolver, requires *model.UserRole) (interface{}, error) {

	// currentUser := auth.GetCurrentUser(ctx)

	// if currentUser == nil {
	// 	return nil, errors.AuthenticationRequired
	// }

	// requiredRole := user.Role(requires.String())

	// if !hasRequiredRole(currentUser.Role, requiredRole) {
	// 	return nil, gqlerror.Errorf(
	// 		"Access denied: requires %s role",
	// 		requiredRole,
	// 	)
	// }

	return next(ctx)
}

func hasRequiredRole(userRole, requiredRole user.Role) bool {

	roleHierarchy := map[user.Role]int{
		user.RoleUSER:  1,
		user.RoleADMIN: 2,
	}

	userLevel := roleHierarchy[userRole]
	if userLevel == 0 {
		userLevel = 1
	}

	requiredLevel := roleHierarchy[requiredRole]
	return userLevel >= requiredLevel
}
