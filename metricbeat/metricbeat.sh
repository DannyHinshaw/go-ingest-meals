#!/usr/bin/env bash
metricbeat modules enable kibana
metricbeat setup
metricbeat -e
