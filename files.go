package leopard

import (
	"errors"
	"fmt"
	"github.com/volix-dev/leopard/files"
	_ "github.com/volix-dev/leopard/files/drivers/osFs"
	_ "github.com/volix-dev/leopard/files/drivers/s3"
)

func getFileDriver() (files.Driver, error) {

	driverName := EnvSettingD("FILE_DRIVER", "local").GetValue().(string)
	fmt.Println(driverName)

	switch driverName {
	case "local":
		return files.Get("os", EnvSettingD("FILE_LOCAL_PATH", "./store").GetValue().(string))
	}

	return nil, errors.New("driver not found")
}
