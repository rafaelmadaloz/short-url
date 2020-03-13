from django.contrib import admin
from .models import Url


class UrlAdmin(admin.ModelAdmin):
    list_display = ('id', 'url', 'short_url', 'last_check', 'last_check_status')

admin.site.register(Url, UrlAdmin)
