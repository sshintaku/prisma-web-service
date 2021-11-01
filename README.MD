This GOLANG based webservice provides a modified REST API service to process data to be used with the Angular reporting front end.  This is currently only two basic REST services but will grow over time.

For this service to work, it needs two exported keys and they are expected to be in caps:

1. APIKEY
2. PASSWORD

In order for the GOLANG service to authenticate agains Prisma Cloud the above two keys need to be exported into the environment via the two above mentioned variables:

e.g.,   export APIKEY="xxxyyyzzz"
        export PASSWORD="aaabbbccc"

