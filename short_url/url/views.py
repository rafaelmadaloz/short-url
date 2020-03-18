from django.shortcuts import redirect, get_object_or_404
from django.urls import reverse_lazy
from django.views import generic

from url.models import Url
from url.forms import UrlForm


def redirect_url(request, hex_id):
    short_url = "http://127.0.0.1:8000/" + hex_id
    url = get_object_or_404(Url, short_url=short_url)
    return redirect(url.url)


class AddUrlView(generic.CreateView):
    model = Url
    template_name = 'url/new.html'
    form_class = UrlForm
    success_url = reverse_lazy('url:url_list')


class UrlListView(generic.ListView):
    model = Url
    template_name = 'url/list.html'
    context_object_name = 'urls'

    def get_queryset(self):
        return Url.objects.all().order_by('-id')
