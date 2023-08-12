package crud

// Dans l'idée un RootController est un controller qui a une liste de controller
// et qui va dispatcher les requetes vers les controllers enfants

// type RootController struct {
// 	controllers []Controller
// }

// Each controller tend to having a model wich some times can only be an action
// So i need to find a way to have a controller that can be a model or an action
// and that can be used in a list or a show or use

// type Show Controller
// type Create Controller
// type Update Controller
// type Controller interface {
// 	// List
// 	// Show
// 	// Create
// 	// Update
// 	// Delete
// 	// Use
// }

// A controller can also be a RootController

// A list is a view that can be used to list a collection of items
// This list is eather Paginated or not
// the list is rendered using a ui model named Table

// type List interface {
// 	// Table
// 	Use
// }

// The Action define the action that can be done on the list
// It can be a Create, Update, Delete, Show, Use
// It's a function that represent the action with the context of it's execution
// The action can be a function that return a
