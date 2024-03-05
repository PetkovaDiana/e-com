from django.contrib import admin
from .models import *
import nested_admin


@admin.register(DeliveryTypeInfo)
class DeliveryTypeInfoAdmin(admin.ModelAdmin):
    list_display = ('title', 'description', 'can_delivery', 'delivery_price',)
    list_display_links = ('title', 'description', 'can_delivery', 'delivery_price',)
    search_fields = ('title',)
    search_help_text = 'Введите название вида доставки'

    def has_delete_permission(self, request, obj=None):
        return False
    

@admin.register(CdekDeliveryInfo)
class CdekDeliveryInfoAdmin(admin.ModelAdmin):
    list_display = ('description',)
    list_display_links = ('description',)

    def has_delete_permission(self, request, obj=None):
        return False


@admin.register(CourierDeliveryInfo)
class CourierDeliveryInfoAdmin(admin.ModelAdmin):
    list_display = ('description', 'courier_delivery_time_info',)
    list_display_links = ('description', 'courier_delivery_time_info',)

    def has_delete_permission(self, request, obj=None):
        return False


@admin.register(CdekDelivery)
class CdekDeliveryAdmin(admin.ModelAdmin):
    list_display = ('delivery_type', 'pick_up_point_address')
    list_display_links = ('delivery_type', 'pick_up_point_address')

    def has_delete_permission(self, request, obj=None):
        return False
    
    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False
    

@admin.register(CourierDelivery)
class CourierDeliveryAdmin(admin.ModelAdmin):
    list_display = ('index', 'address', 'entrance', 'floor', 'apartment_office', 'delivery_type')
    list_display_links = ('index', 'address', 'entrance', 'floor', 'apartment_office', 'delivery_type')
    list_filter = ('delivery_type',)
    list_per_page = 50
    search_fields = ('index', 'address',)
    search_help_text = 'Введите индекс или адрес'

    def has_delete_permission(self, request, obj=None):
        return False
    
    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False


@admin.register(SelfDelivery)
class SelfDeliveryAdmin(admin.ModelAdmin):
    list_display = ('delivery_type', 'pick_up_point',)
    list_display_links = ('delivery_type', 'pick_up_point',)
    list_filter = ('pick_up_point',)
    list_per_page = 50

    def has_delete_permission(self, request, obj=None):
        return False
    
    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False
    

@admin.register(CourierDeliveryTimeInfo)
class CourierDeliveryTimeInfoAdmin(admin.ModelAdmin):
    list_display = ('mon', 'tue', "wen", "thu", "fri", "sat", "sun",)
    list_display_links = ('mon', 'tue', "wen", "thu", "fri", "sat", "sun",)

    def has_delete_permission(self, request, obj=None):
        return False
    

class CourierDeliveryInline(nested_admin.NestedStackedInline):
    model = CourierDelivery
    extra = 100
    verbose_name = 'Курьерская доставка'
    verbose_name_plural = 'Курьерские доставки'


class CdekDeliveryInline(nested_admin.NestedStackedInline):
    model = CdekDelivery
    extra = 100
    verbose_name = 'Доставка CDEK'
    verbose_name_plural = 'Доставка CDEK'


class SelfDeliveryInline(nested_admin.NestedStackedInline):
    model = SelfDelivery
    extra = 100
    verbose_name = 'Самовывоз'
    verbose_name_plural = 'Самовывозы'
    

@admin.register(DeliveryType)
class DeliveryTypeAdmin(admin.ModelAdmin):
    list_display = ('id', 'order',)
    list_display_links = ('id', 'order',)
    search_fields = ('id',)
    search_help_text = 'Введите id заказа'
    inlines = [CourierDeliveryInline, CdekDeliveryInline, SelfDeliveryInline]

    def has_delete_permission(self, request, obj=None):
        return False
    
    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False