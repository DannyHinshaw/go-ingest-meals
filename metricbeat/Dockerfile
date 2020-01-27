FROM docker.elastic.co/beats/metricbeat:7.5.2

COPY metricbeat.yml /usr/share/metricbeat/metricbeat.yml
COPY wait-for-it.sh .
COPY metricbeat.sh .

USER root
RUN mkdir /var/log/metricbeat \
    && chown metricbeat /usr/share/metricbeat/metricbeat.yml \
    && chmod go-w /usr/share/metricbeat/metricbeat.yml \
    && chown metricbeat /var/log/metricbeat \
    && chmod +x wait-for-it.sh \
    && chmod +x metricbeat.sh

USER metricbeat
