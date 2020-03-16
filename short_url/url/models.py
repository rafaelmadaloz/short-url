from django.db import models


class Url(models.Model):
    url = models.CharField(max_length=512)
    short_url = models.CharField(max_length=64, blank=True)
    last_check = models.DateTimeField(blank=True, null=True)
    last_check_status = models.IntegerField(blank=True, null=True)

    def __str__(self):
        return self.url

    def save(self):
        if not self.pk:
            super().save()

        self.short_url = "http://127.0.0.1:8000/" + str(hex(self.pk)[2:])
        return super().save()
