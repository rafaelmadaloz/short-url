from django.urls import path

from . import views

app_name = 'url'
urlpatterns = [
    path('<str:hash_id>/', views.redirect_url, name='redirect'),
    path('', views.UrlListView.as_view(), name='url_list'),
    path('add', views.AddUrlView.as_view(), name='add_url'),
    path('delete/<str:hash_id>', views.delete_url, name='delete_url'),
]
