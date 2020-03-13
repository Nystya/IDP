from django.urls import path
from django.urls import include
from django.conf.urls import url
from django.views.generic import TemplateView
from allauth.account.views import email, logout

urlpatterns = [
    # this url is used to generate email content
    url(r'^password-reset/confirm/(?P<uidb64>[0-9A-Za-z_\-]+)/(?P<token>[0-9A-Za-z]{1,13}-[0-9A-Za-z]{1,20})/$',
        TemplateView.as_view(template_name="password_reset_confirm.html"),
        name='password_reset_confirm'),

    path('',  include('rest_auth.urls')),
    path('registration/', include('rest_auth.registration.urls')),
    url(r"^email/$", email, name="account_email"),
    url(r"^logout/$", email, name="account_logout"),
]