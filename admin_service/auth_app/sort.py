from django.contrib import admin

class UserTypeFilter(admin.SimpleListFilter):
    title = 'Тип пользователя'
    parameter_name = 'user_type'

    def lookups(self, request, model_admin):
        return (
            ('phys', 'Физ. лицо'),
            ('legal', 'Юр. лицо'),
        )

    def queryset(self, request, queryset):
        if self.value() == 'phys':
            return queryset.filter(inn__exact='')
        elif self.value() == 'legal':
            return queryset.exclude(inn__exact='')