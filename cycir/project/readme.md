# Instructions
## Host modification
Add backend entry in hosts file, for windows, at C:\Windows\System32\drivers\etc

![2022-11-08_10-35-25](https://user-images.githubusercontent.com/30485720/200469402-0d1e55f9-4734-4f0f-a992-bdb7296ece3d.png)


## Installation
### Docker
~~~
cd projects
docker compose up
~~~

### docker swarm 
~~~
cd projects
// init
docker swarm init
docker stack deploy -c swarm.yml cycir

// scale services
docker service ls
docker service scale {service-name}=3
docker service rm {service-name}

// stop
docker stack rm cycir
docker swarm leave --force
~~~
> **For docker swarm, after deploying, remove project_migrate-service (this is used for migrations)**
 
