package assets

import "github.com/nyaruka/goflow/envs"

// LocationHierarchy is a searchable hierarchy of locations.
//
//   {
//     "name": "Rwanda",
//     "aliases": ["Ruanda"],
//     "children": [
//       {
//         "name": "Kigali City",
//         "aliases": ["Kigali", "Kigari"],
//         "children": [
//           {
//             "name": "Gasabo",
//             "children": [
//               {
//                 "id": "575743222",
//                 "name": "Gisozi"
//               },
//               {
//                 "id": "457378732",
//                 "name": "Ndera"
//               }
//             ]
//           },
//           {
//             "name": "Nyarugenge",
//             "children": []
//           }
//         ]
//       },
//       {
//         "name": "Eastern Province"
//       }
//     ]
//   }
//
// @asset location
type LocationHierarchy interface {
	FindByPath(path envs.LocationPath) *envs.Location
	FindByName(name string, level envs.LocationLevel, parent *envs.Location) []*envs.Location
}
