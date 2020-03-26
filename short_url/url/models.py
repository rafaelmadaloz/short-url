from django.db import models

from url.utils.hash import hashids


class Url(models.Model):
    url = models.URLField(max_length=512, unique=True)
    last_check = models.DateTimeField(blank=True, null=True)
    last_check_status = models.IntegerField(blank=True, null=True)

    @property
    def short_url(self):
        return hashids.encode(self.pk)

    def __str__(self):
        return self.url
