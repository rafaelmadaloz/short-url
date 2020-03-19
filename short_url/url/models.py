from django.db import models

from hashids import Hashids


class Url(models.Model):
    url = models.URLField(max_length=512, unique=True)
    short_url = models.URLField(max_length=64, blank=True, unique=True)
    last_check = models.DateTimeField(blank=True, null=True)
    last_check_status = models.IntegerField(blank=True, null=True)

    def __str__(self):
        return self.url

    def save(self):
        if not self.pk:
            super().save()

        hashids = Hashids(salt='IIm54tostyz6tWoIJukG')
        hashid = hashids.encode(self.pk)

        self.short_url = hashid
        return super().save()
