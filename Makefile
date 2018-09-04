build:
	docker build --no-cache -t quota-setter:0.1 .
	docker tag quota-setter:0.1 ai-image.jd.com/datamgr/quota-setter:0.1
	docker push ai-image.jd.com/datamgr/quota-setter:0.1
deploy-only:
	kubectl create ns cephquota || true
	kubectl delete -f deployment || true
	kubectl apply -f deployment
deploy: build deploy-only