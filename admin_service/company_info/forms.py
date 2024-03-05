from ckeditor_uploader.widgets import CKEditorUploadingWidget
from django import forms
from django.contrib import admin

from .models import Requisites, PrivacyPolicy


class RequisitesForm(forms.ModelForm):
    text = forms.CharField(widget=CKEditorUploadingWidget(), label='Описание')

    class Meta:
        model = Requisites
        fields = '__all__'


class RequisitesAdmin(admin.ModelAdmin):
    form = RequisitesForm


class PrivacyPolicyForm(forms.ModelForm):
    text = forms.CharField(widget=CKEditorUploadingWidget(), label='Описание')

    class Meta:
        model = PrivacyPolicy
        fields = '__all__'


class PrivacyPolicyAdmin(admin.ModelAdmin):
    form = PrivacyPolicyForm