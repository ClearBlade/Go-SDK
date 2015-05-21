package GoSDK

//This file provides the interface for establishing connect collections
//that is to say, collections that are interfaced with non-platform databases
//they have to be treated a little bit differently, because a lot of configuration information
//needs to be trucked across the line during setup. enough that it's more helpful to have it in a
//struct than it is just in a map, or an endless list of function arguments.
