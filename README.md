# Download counter

This is work in progress.

The intent is to push download counts into a timeseries database (InfluxDB) based on nginx or apache logfile(s), or any file matching a regex.

Tags are added to the timeseries entry, based on fields defined in the regex. Note that named fields are not supported, since additional features might be added in the future that would prevent this. 



