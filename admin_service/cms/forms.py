from ckeditor_uploader.widgets import CKEditorUploadingWidget
from django import forms
from django.contrib import admin

from .models import CurrentPromotions, Blog


class CurrentPromotionsForm(forms.ModelForm):
    description = forms.CharField(widget=CKEditorUploadingWidget(), label='Описание')

    class Meta:
        model = CurrentPromotions
        fields = '__all__'


class CurrentPromotionsAdmin(admin.ModelAdmin):
    form = CurrentPromotionsForm


class BlogForm(forms.ModelForm):
    description = forms.CharField(widget=CKEditorUploadingWidget(), label='Описание')

    class Meta:
        model = Blog
        fields = '__all__'


class BlogAdmin(admin.ModelAdmin):
    form = BlogForm