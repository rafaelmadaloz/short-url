from django.urls import path

from . import views

app_name = 'url'
urlpatterns = [
    path('<str:hex_id>/', views.redirect_url, name='redirect'),
    path('', views.UrlListView.as_view(), name='url_list'),
    path('new', views.AddUrlView.as_view(), name='add_url'),
]
