build:
	docker build --no-cache -t cephfs-quota:0.1 .
	docker tag quota-setter:0.1 ai-image.jd.com/ceph/cephfs-quota:0.1
	docker push ai-image.jd.com/ceph/cephfs-quota:0.1
deploy-only:
	kubectl create ns cephquota || true
	kubectl delete -f deployment || true
	kubectl apply -f deployment
clean:
	kubectl delete ns cephquota || true
deploy: build deploy-only

