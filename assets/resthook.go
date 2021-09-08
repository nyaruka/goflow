package assets

// Resthook is a set of URLs which are subscribed to the named event.
//
//   {
//     "slug": "new-registration",
//     "subscribers": [
//       "http://example.com/record.php?@contact.uuid"
//     ]
//   }
//
// @asset resthook
type Resthook interface {
	Slug() string
	Subscribers() []string
}
