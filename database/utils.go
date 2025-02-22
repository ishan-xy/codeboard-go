package database
import(
	"context"
	"log"
	"errors"
	utils "github.com/ItsMeSamey/go_utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)
func (c *Collection[T]) GetExists(filter any) (out T, exists bool, err error) {
	result := c.FindOne(context.Background(), filter)
	err = result.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return out, false, nil
		}
		log.Printf("Error finding document: %v\n", err)
		return out, false, utils.WithStack(err)
	}

	if err := result.Decode(&out); err != nil {
		log.Printf("Error decoding document: %v\n", err)
		return out, false, utils.WithStack(err)
	}

	return out, true, nil
}