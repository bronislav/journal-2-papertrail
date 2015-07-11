FROM       scratch
MAINTAINER Anton Ilin <anton@ilin.dn.ua>
ADD        journal-2-papertrail journal-2-papertrail
ENTRYPOINT ["/journal-2-papertrail"]
