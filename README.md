A minimal template app for go backends that require:

## Rationale:

In most project the hardest part for me has been the scaffolding.  Structuring the protos, structuring the web pages, setting up the binaries.  Id have fun building the core components and id be forced to look at the web aspect and it is a huge sync - from design, to layouts, to tool selection etc.  This repo just gives a couple of blueprints I can use for 95% of my app idea use cases.  Just global search and replace a few variables and you have a working app.  You just have to start adding files in the respective places to "grow" the functionality instead of scrambling to figure out web conventions etc.  One of the biggest advantages is with vibecoding.  We can now just have a running up that can direct the agent to to query etc instead having to figure out how to even create a web seutp which each time comes up with a different result!

## Global Search and Replace:

weewar -> Your backend name, eg "MyFancyGame"
weewar -> module name for your backend, age "myfancygame"
AppItem -> This is one of the entities in our proto, you can add more or replace this with the name of your key entity.  Feel free to add more entities that make sense for your system
appitem -> "Variable" names for the AppItem type in the code

./services/appitems_service.go - Definition of your service for appitems.  You can one one file for each entity type

## Key technologies and Stack components:

1. Protos/GRPC services
2. API Fronted by a gateway service
3. Powered by OneAuth for oauth
4. Basic frontend based on Templar go-templates - this can be customized in the future.
5. Tailwind for styling and Typescript for front end.  Add react/vue etc if youd like
6. Webpack for complex pages
7. Many sample page templates, eg ListPages, BorderLayout pages, pagse with DockView etc you can easily copy for other
   parts of your page.

Other backend choices (like datastores) are upto the app/service dev, eg:

1. Which services are to be added.
2. Which backends are to be used (for storage, etc)
3. How to deploy them to specific hosting providers (eg appengine, heroku etc)
4. Selecting frontend frameworks.

## Requirements

1. Basic/Standard Go Tooling:

* Go
* Air (for fast reloads)
* Protobuf
* GRPC
* Buf (for generating artificates from grpc protos)
* Webpack for any complex pages

## Getting Started

1. Clone this Repo

Replace the following variables:

TODO:
1. Common docker compose manifests for packaging for development.
2. Optional k8s configs if needed in the future for testing against cluster deployments

## Conventions:

1. Protos are defined in ./protos folder.  Grpc is our source of truth for everything.  Every other client is generated
   with this.   Connect clients, gateway bindings, TS clients even MCP tool names!
2. Web server, handlers, templates are all defined in the web folder
3. services folder takes care of the core proto service implementation and all the backend heavy lifting
4. The main.go just loads an "App" type that can run a bunch of servers (GrpcServer, WebServer and any others).  Each
   "cluster" has its own server - services has the GrpcServer, web has the WebServer (with grpc bindings).
