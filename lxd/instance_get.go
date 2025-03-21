package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lxc/lxd/lxd/instance"
	"github.com/lxc/lxd/lxd/response"
)

// swagger:operation GET /1.0/instances/{name} instances instance_get
//
// Get the instance
//
// Gets a specific instance (basic struct).
//
// ---
// produces:
//   - application/json
// parameters:
//   - in: query
//     name: project
//     description: Project name
//     type: string
//     example: default
// responses:
//   "200":
//     description: Instance
//     schema:
//       type: object
//       description: Sync response
//       properties:
//         type:
//           type: string
//           description: Response type
//           example: sync
//         status:
//           type: string
//           description: Status description
//           example: Success
//         status_code:
//           type: integer
//           description: Status code
//           example: 200
//         metadata:
//           $ref: "#/definitions/Instance"
//   "403":
//     $ref: "#/responses/Forbidden"
//   "500":
//     $ref: "#/responses/InternalServerError"

// swagger:operation GET /1.0/instances/{name}?recursion=1 instances instance_get_recursion1
//
// Get the instance
//
// Gets a specific instance (full struct).
//
// recursion=1 also includes information about state, snapshots and backups.
//
// ---
// produces:
//   - application/json
// parameters:
//   - in: query
//     name: project
//     description: Project name
//     type: string
//     example: default
// responses:
//   "200":
//     description: Instance
//     schema:
//       type: object
//       description: Sync response
//       properties:
//         type:
//           type: string
//           description: Response type
//           example: sync
//         status:
//           type: string
//           description: Status description
//           example: Success
//         status_code:
//           type: integer
//           description: Status code
//           example: 200
//         metadata:
//           $ref: "#/definitions/Instance"
//   "403":
//     $ref: "#/responses/Forbidden"
//   "500":
//     $ref: "#/responses/InternalServerError"
func instanceGet(d *Daemon, r *http.Request) response.Response {
	instanceType, err := urlInstanceTypeDetect(r)
	if err != nil {
		return response.SmartError(err)
	}

	projectName := projectParam(r)
	name := mux.Vars(r)["name"]

	// Parse the recursion field
	recursionStr := r.FormValue("recursion")

	recursion, err := strconv.Atoi(recursionStr)
	if err != nil {
		recursion = 0
	}

	// Handle requests targeted to a container on a different node
	resp, err := forwardedResponseIfInstanceIsRemote(d, r, projectName, name, instanceType)
	if err != nil {
		return response.SmartError(err)
	}
	if resp != nil {
		return resp
	}

	c, err := instance.LoadByProjectAndName(d.State(), projectName, name)
	if err != nil {
		return response.SmartError(err)
	}

	var state interface{}
	var etag interface{}
	if recursion == 0 {
		state, etag, err = c.Render()
	} else {
		state, etag, err = c.RenderFull()
	}
	if err != nil {
		return response.SmartError(err)
	}

	return response.SyncResponseETag(true, state, etag)
}
