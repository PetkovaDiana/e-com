from django.db.models import Q
from django.utils.translation import gettext_lazy as _
from django.contrib import admin
from django.db.models import Count


class RatingFilter(admin.SimpleListFilter):
    title = _('рейтинг')
    parameter_name = 'rating'

    def lookups(self, request, model_admin):
        return (
            ('0-1', _('0-1')),
            ('1-2', _('1-2')),
            ('2-3', _('2-3')),
            ('3-4', _('3-4')),
            ('4-5', _('4-5')),
        )

    def queryset(self, request, queryset):
        if self.value() == '0-1':
            return queryset.filter(Q(rating__gte=0) & Q(rating__lt=1))
        elif self.value() == '1-2':
            return queryset.filter(Q(rating__gte=1) & Q(rating__lt=2))
        elif self.value() == '2-3':
            return queryset.filter(Q(rating__gte=2) & Q(rating__lt=3))
        elif self.value() == '3-4':
            return queryset.filter(Q(rating__gte=3) & Q(rating__lt=4))
        elif self.value() == '4-5':
            return queryset.filter(Q(rating__gte=4) & Q(rating__lte=5))
        else:
            return queryset
        

class SalesCountFilter(admin.SimpleListFilter):
    title = _('кол-во продаж')
    parameter_name = 'sales_count'

    def lookups(self, request, model_admin):
        return (
            ('0-50', _('0-50')),
            ('50-100', _('50-100')),
            ('100-150', _('100-150')),
            ('150-200', _('150-200')),
            ('200-250', _('200-250')),
            ('250-300', _('250-300')),
            ('300-350', _('300-350')),
            ('350-500', _('350-500')),
            ('500-1000', _('500-1000')),
            ('1000-1500', _('1000-1500')),
            ('1500-3000', _('1500-3000')),
            ('over_3000', _('>3000')),
        )

    def queryset(self, request, queryset):
        if self.value() == '0-50':
            return queryset.filter(sales_count__range=(0, 50))
        elif self.value() == '50-100':
            return queryset.filter(sales_count__range=(50, 100))
        elif self.value() == '100-150':
            return queryset.filter(sales_count__range=(100, 150))
        elif self.value() == '150-200':
            return queryset.filter(sales_count__range=(150, 200))
        elif self.value() == '200-250':
            return queryset.filter(sales_count__range=(200, 250))
        elif self.value() == '250-300':
            return queryset.filter(sales_count__range=(250, 300))
        elif self.value() == '300-350':
            return queryset.filter(sales_count__range=(300, 350))
        elif self.value() == '350-500':
            return queryset.filter(sales_count__range=(350, 500))
        elif self.value() == '500-1000':
            return queryset.filter(sales_count__range=(500, 1000))
        elif self.value() == '1000-1500':
            return queryset.filter(sales_count__range=(1000, 1500))
        elif self.value() == '1500-3000':
            return queryset.filter(sales_count__range=(1500, 3000))
        elif self.value() == 'over_3000':
            return queryset.filter(sales_count__gte=3000)


class RequestDetailCountFilter(admin.SimpleListFilter):
    title = _('кол-во детальных просмотров')
    parameter_name = 'request_detail_count'

    def lookups(self, request, model_admin):
        return (
            ('0-50', _('0-50')),
            ('50-100', _('50-100')),
            ('100-150', _('100-150')),
            ('150-200', _('150-200')),
            ('200-250', _('200-250')),
            ('250-300', _('250-300')),
            ('300-350', _('300-350')),
            ('350-500', _('350-500')),
            ('500-1000', _('500-1000')),
            ('1000-1500', _('1000-1500')),
            ('1500-3000', _('1500-3000')),
            ('over_3000', _('>3000')),
        )

    def queryset(self, request, queryset):
        if self.value() == '0-50':
            return queryset.filter(sales_count__range=(0, 50))
        elif self.value() == '50-100':
            return queryset.filter(sales_count__range=(50, 100))
        elif self.value() == '100-150':
            return queryset.filter(sales_count__range=(100, 150))
        elif self.value() == '150-200':
            return queryset.filter(sales_count__range=(150, 200))
        elif self.value() == '200-250':
            return queryset.filter(sales_count__range=(200, 250))
        elif self.value() == '250-300':
            return queryset.filter(sales_count__range=(250, 300))
        elif self.value() == '300-350':
            return queryset.filter(sales_count__range=(300, 350))
        elif self.value() == '350-500':
            return queryset.filter(sales_count__range=(350, 500))
        elif self.value() == '500-1000':
            return queryset.filter(sales_count__range=(500, 1000))
        elif self.value() == '1000-1500':
            return queryset.filter(sales_count__range=(1000, 1500))
        elif self.value() == '1500-3000':
            return queryset.filter(sales_count__range=(1500, 3000))
        elif self.value() == 'over_3000':
            return queryset.filter(sales_count__gte=3000)


class DiscountPercentFilter(admin.SimpleListFilter):
    title = _('Вид скидки')
    parameter_name = 'discount_percent'

    def lookups(self, request, model_admin):
        return (
            ('not_set', _('В процентах %')),
            ('set', _('В рублях')),
        )

    def queryset(self, request, queryset):
        if self.value() == 'set':
            return queryset.annotate(discount_percent_count=Count('discount_percent')).filter(discount_percent_count=0)
        elif self.value() == 'not_set':
            return queryset.exclude(discount_percent__isnull=True)
