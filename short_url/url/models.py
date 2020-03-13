from django.db import models

class Url(models.Model):
    url = models.CharField(max_length=512)
    short_url = models.CharField(max_length=64)
    last_check = models.DateTimeField(blank=True, null=True)
    last_check_status = models.IntegerField(blank=True, null=True)

    def __str__(self):
        return self.url