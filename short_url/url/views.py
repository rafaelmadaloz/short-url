from django.shortcuts import redirect, get_object_or_404

from .models import Url


def redirect_url(request, url_id):
    url = get_object_or_404(Url, pk=url_id)
    return redirect(url.url)
