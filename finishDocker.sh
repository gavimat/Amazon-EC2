echo "Finalizando todos os containers Docker"
docker kill $(docker ps -q)

echo "Removendo todos os containers Docker"
docker rm -f `docker ps -aq`

echo "Deletando todos os volumes Docker"
docker volume rm $(docker volume ls -q)