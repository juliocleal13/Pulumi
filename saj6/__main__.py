import pulumi
import pulumi_docker as docker
from pulumi_docker import ImageRegistry, RemoteImage, Container

config = pulumi.Config()
username = config.require('dockerUsername')
password = config.require_secret('dockerPassword')
app = 'unj-workflow-frontend'
registry_url = 'docker-unj-repo.softplan.com.br/unj'

tags = config.require_object("tags")

registry_image = docker.get_registry_image(name=f'{registry_url}/{app}:{tags.get("unj-workflow-frontend")}')

remote_image = docker.RemoteImage(app,name=registry_image.name,
    pull_triggers=[registry_image]
)