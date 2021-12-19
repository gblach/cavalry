FROM almalinux:minimal

RUN microdnf upgrade
RUN microdnf install python39
RUN microdnf clean all
RUN pip3 install -U pip
RUN pip3 install falcon gunicorn
RUN pip3 cache purge

COPY example.py .
CMD ["gunicorn", "--bind=0.0.0.0:8000", "example:app"]

EXPOSE 8000/tcp
