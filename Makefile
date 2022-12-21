
run_in_mem:
	docker-compose build app_in_mem
	docker-compose up  app_in_mem

run_in_db:
	docker-compose build  app_in_db
	docker-compose up  app_in_db

.PHONY: run_in_db, run_in_mem