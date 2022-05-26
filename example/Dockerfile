FROM almalinux/9-minimal

RUN microdnf upgrade -y \
	&& microdnf install -y python3-pip \
	&& microdnf clean all

RUN pip3 install --upgrade pip \
	&& pip3 install falcon gunicorn \
	&& pip3 cache purge

COPY example.py .
CMD ["gunicorn", "--bind=0.0.0.0:8000", "example:app"]

EXPOSE 8000/tcp
