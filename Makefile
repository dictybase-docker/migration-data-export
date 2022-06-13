remove:
	rm -rf data
docker-copy: remove
	docker cp dsc-annotations:/data/ $(shell pwd)/ 
	docker cp dictybase-users:/data/ $(shell pwd)/
	docker cp literature:/data/ $(shell pwd)/
	docker cp genome-annotations:/data/ $(shell pwd)/
create-tarball: docker-copy 
		cd data/ \
		&& tar cvzf stockcenter.tar.gz stockcenter \
		&& tar cvzf users.tar.gz users \ 
		&& tar cvzf genome-annotation.tar.gz annotation \
		&& tar cvzf literature.tar.gz literature
