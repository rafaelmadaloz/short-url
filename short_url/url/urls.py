from django.urls import path

from . import views

app_name = 'url'
urlpatterns = [
    path('<int:url_id>', views.redirect_url, name='redirect'),
]
