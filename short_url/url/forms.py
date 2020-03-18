from django import forms

from url.models import Url

class UrlForm(forms.ModelForm):
    url = forms.CharField(
                        label="URL",
                        max_length=512,
                        widget=forms.TextInput(attrs={'class': 'form-control'}))
    class Meta:
        model = Url
        fields = ['url']
