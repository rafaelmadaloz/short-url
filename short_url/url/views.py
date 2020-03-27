from django.shortcuts import redirect, get_object_or_404
from django.urls import reverse_lazy
from django.views import generic
from django.utils.translation import gettext as _

from hashids import Hashids

from url.models import Url
from url.forms import UrlForm
from url.utils.hash import hashids


def redirect_url(request, hash_id):
    hashid = hashids.decode(hash_id)
    url_id = str(hashid[0])
    url = get_object_or_404(Url, pk=url_id)
    return redirect(url.url)


def delete_url(request, hash_id):
    hashid = hashids.decode(hash_id)
    url_id = str(hashid[0])
    url = Url.objects.filter(pk=url_id).delete()
    return redirect('url:url_list')


class AddUrlView(generic.CreateView):
    model = Url
    template_name = 'url/add.html'
    form_class = UrlForm
    success_url = reverse_lazy('url:url_list')

    def get_context_data(self, **kwargs):
        data = super().get_context_data(**kwargs)
        data['submit_button'] = _('Submit')

        return data


class UrlListView(generic.ListView):
    model = Url
    template_name = 'url/list.html'
    context_object_name = 'urls'

    def get_context_data(self, **kwargs):
        data = super().get_context_data(**kwargs)
        data['table_indexes'] = {
            'short_url': _('Short URL'),
            'url': _('URL'),
            'last_check': _('Last Check'),
            'status': _('Status'),
            'delete': _('Delete'),
        }
        data['add_button'] = _('Add')

        return data

    def get_queryset(self):
        return Url.objects.all().order_by('-id')
