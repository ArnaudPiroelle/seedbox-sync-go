FROM scratch

COPY build/seedbox-sync-Linux-x86_64 /opt/seedbox-sync

VOLUME /config
ENTRYPOINT [ "/opt/seedbox-sync" ]
CMD [ "scheduler", "-c", "/config/config.json"]