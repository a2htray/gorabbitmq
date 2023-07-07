# enable plugins
rabbitmq-plugins enable rabbitmq_management

rabbitmqctl add_user goadmin 123456
rabbitmqctl set_user_tags goadmin administrator
rabbitmqctl add_vhost goapp-vhost --description 'virtual host for go test application' --tags goapp
rabbitmqctl set_permissions -p goapp-vhost goadmin ".*" ".*" ".*"


