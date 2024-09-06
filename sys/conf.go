package sys

import (
	"github.com/spf13/viper"
	"log"
)

/*改进建议：
- 配置文件路径的更灵活处理：
  目前你在 init 函数中使用了硬编码的配置文件路径，即 "inputs" 和 "../inputs"。这样的路径可能不够灵活，特别是当你的应用程序的目录结构发生变化时。
  考虑使用环境变量、命令行参数或者配置文件来指定配置文件的路径。这样，你的应用程序可以更容易地适应不同的部署环境。
- 支持多个配置文件：
  如果你的应用程序需要加载多个配置文件（例如，一个通用配置文件和一个特定环境的配置文件），可以考虑支持多个配置文件。
  viper 可以通过 SetConfigName 设置多个配置文件名，然后按照优先级顺序加载它们。
- 错误处理的改进：
  目前你的错误处理是通过日志记录的，但这可能不足够灵活。你可以考虑返回错误对象，以便调用方可以更好地处理错误。
  另外，你可以在 ReadConfig 函数中添加更多的错误处理逻辑，例如处理配置文件格式错误等。
*/

var Conf *viper.Viper

func init() {
	Conf = viper.New()
	Conf.AddConfigPath("inputs")
	Conf.AddConfigPath("../inputs") // for apps at the same level as inputs
	Conf.SetConfigName("config")
	Conf.SetConfigType("yaml")
	ReadConfig(Conf)
}

func ReadConfig(c *viper.Viper) {
	if err := c.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found")
		} else {
			log.Println("Config file found but error reading:", err)
		}
	}
}
