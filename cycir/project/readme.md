# Instructions
## Login account

account: admin@example.com

password: password

## Local Settings
Add backend entry in hosts file, for windows, at C:\Windows\System32\drivers\etc

![2022-11-08_10-35-25](https://user-images.githubusercontent.com/30485720/200469402-0d1e55f9-4734-4f0f-a992-bdb7296ece3d.png)

## App Settings
### Overview -> Settings
SMTP server: mailhog (instead of localhost)

SMTP port: 1025


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
> **default postgres port is on 5432, change it if your local postgres already resides at port 5432**
 
## Potential bugs
> **After logout, please press F5 before login in again, for some reasons, bootstrap form doesn't submit form data**

> **If in any cases, the app breaks, please open devtools, removes all caches and start again (this could happen because of the reason mentioned above, very sorry about that)**

