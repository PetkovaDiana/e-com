import requests
import csv

from django.contrib import admin
from django.db.models import Q
from django import forms
from django.http import HttpResponse
from .sort import *

from django.template.response import TemplateResponse
from django.db.models import Avg


from .models import *

def get_average_rating():
    return SiteReview.objects.aggregate(Avg('rating'))['rating__avg']

@admin.register(User)
class UserAdmin(admin.ModelAdmin):
    list_display = ('id', 'name', 'surname', 'phone', 'inn', 'company_name', 'company_address', 'manager_name', 'kpp')
    list_display_links = ('id', 'name', 'surname', 'phone', 'inn', 'company_name')
    list_per_page = 30
    search_fields = ('id', 'name', 'phone', 'company_name')
    search_help_text = 'Введите id, имя, телефон или имя компании'
    list_filter = (UserTypeFilter,)


    def get_queryset(self, request):
        query = super(UserAdmin, self).get_queryset(request)
        filtered_query = query.filter(~Q(phone=''))
        return filtered_query

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False


class UnregisteredUserAdmin(UserAdmin):
    list_display = ('id',)
    list_display_links = ('id',)
    list_per_page = 30
    search_fields = ('id',)
    search_help_text = 'Введите id'

    class Meta:
        proxy = True
        verbose_name = 'Незарегестрированный пользователь'
        verbose_name_plural = 'Незарегестрированные пользователи'

    def get_queryset(self, request):
        return self.model.objects.filter(phone='')

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False


@admin.action(description='Сделать рассылку выбранных категорий')
def sender(modeladmin, request, queryset):
    emails_body = list(queryset.values_list('body', flat=True))
    emails_id_dict = list(queryset.values_list('emails_array', flat=True))
    post_data = {
        "key": "lydluxlhxlhx637447JfjfHzhg",
        "emails": list(emails_id_dict),
        "body": emails_body
    }
    response = requests.post('http://127.0.0.1:8080/email_sender', json=post_data)


class EmailListCsvsAdmin(forms.Form):
    csv_upload = forms.FileField()


@admin.register(Email)
class EmailListToSendAdmin(admin.ModelAdmin):
    list_display = ('id', 'email', 'can_to_send_news', 'can_to_send_personal_offers', 'user',)
    list_display_links = ('id', 'email', 'can_to_send_news', 'can_to_send_personal_offers', 'user',)
    list_per_page = 30
    search_fields = ('id', 'email',)
    search_help_text = 'Введите id или email'

    def export_csv(self, request, queryset):
        response = HttpResponse(content_type='text/csv')
        response['Content-Disposition'] = 'attachment; filename="emails.csv"'
        writer = csv.writer(response)

        writer.writerow(['email'])
        for email in queryset:
            writer.writerow([email.email])
        return response

    export_csv.short_description = "Экспортировать выбранные email в формате CSV"

    actions = [export_csv]

def create_modeladmin(modeladmin, model, name = None, data = None):

    class Meta:
        proxy = data.proxy
        app_label = model._meta.app_label
        verbose_name = data.verbose_name
        verbose_name_plural = data.verbose_name_plural

    attrs = {'__module__': '', 'Meta': Meta}

    newmodel = type(Meta.verbose_name_plural, (model,), attrs)

    admin.site.register(newmodel, modeladmin)
    return modeladmin

create_modeladmin(UnregisteredUserAdmin, model=User, data=UnregisteredUserAdmin.Meta)


admin.site.site_url = "https://ufaelectro.ru/"

@admin.register(SiteReview)
class SiteReviewAdmin(admin.ModelAdmin):
    list_display = ('rating', 'comment', 'created_at',)
    list_display_links = ('rating', 'comment', 'created_at',)
    list_per_page = 30
    search_fields = ('rating', 'comment',)
    search_help_text = 'Введите рейтинг или комментарий'

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False

    def get_average_rating(self):
        return get_average_rating()

    def changelist_view(self, request, extra_context=None):
        if extra_context is None:
            extra_context = {}

        extra_context['average_rating'] = self.get_average_rating()
        return super().changelist_view(request, extra_context=extra_context)