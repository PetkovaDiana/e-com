from .models import *
from django.contrib import admin
from django.db.models import Q
from django.utils.translation import gettext_lazy as _
from .sort import *
from django.urls import reverse
from django.utils.html import format_html
from delivery.models import DeliveryType
from delivery.admin import CdekDeliveryInline, CourierDeliveryInline, SelfDeliveryInline
import nested_admin

class CartProductInline(admin.TabularInline):
    model = Cart.cart_product.through
    extra = 1
    verbose_name = 'Продукт'
    verbose_name_plural = 'Продукты'


class FilesInline(admin.TabularInline):
    model = ProductFiles.products.through
    extra = 1
    verbose_name = 'Продукт'
    verbose_name_plural = 'Продукты'


@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ('uuid', 'title', 'can_to_view', 'level',)
    list_display_links = ('uuid', 'title', 'can_to_view', 'level',)
    list_per_page = 50
    list_filter = ('level', 'can_to_view',)
    search_fields = ('uuid', 'title',)
    search_help_text = 'Введите id или имя категории'

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ('uuid', 'title', 'base_unit', 'price', 'vendor_code',)
    list_display_links = ('uuid', 'title', 'base_unit', 'price', 'vendor_code',)
    list_per_page = 50
    list_filter = ('can_to_view',)
    search_fields = ('uuid', 'title', 'vendor_code',)
    search_help_text = 'Введите id, имя или вендор код продукта'

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False



@admin.register(ProductStatistic)
class ProductStatisticAdmin(admin.ModelAdmin):
    list_display = ('id', 'rating', 'sales_count', 'request_detail_count', 'product_uuid',)
    list_display_links = ('id', 'rating', 'sales_count', 'request_detail_count', 'product_uuid',)
    list_per_page = 50
    search_fields = ('uuid', 'title', 'vendor_code',)
    search_help_text = 'Введите id, имя или вендор код продукта'
    list_filter = (RatingFilter, SalesCountFilter, RequestDetailCountFilter,)

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False

    def has_change_permission(self, request, obj=None):
        return False


@admin.register(PromoCode)
class PromoCodeAdmin(admin.ModelAdmin):
    list_display = ('id', 'promo_code', 'discount_percent', 'discount_sum', 'number_of_uses', 'expires_at',)
    list_display_links = ('id', 'promo_code', 'discount_percent', 'discount_sum', 'number_of_uses', 'expires_at',)
    list_filter = (DiscountPercentFilter,)
    list_per_page = 50
    search_fields = ('id', 'promo_code',)
    search_help_text = 'Введите id или промокод'


@admin.register(CartProduct)
class CartProductAdmin(admin.ModelAdmin):
    list_display = ('id', 'count', 'total_price',)
    list_display_links = ('id', 'count', 'total_price',)

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False


@admin.register(Cart)
class CartAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'in_order',)
    list_display_links = ('id', 'user', 'in_order',)
    list_per_page = 50
    inlines = [CartProductInline, ]
    list_filter = ('in_order',)
    search_fields = ('id',)
    search_help_text = 'Введите id корзины'
    
    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False


@admin.register(Order)
class OrderAdmin(admin.ModelAdmin):
    list_display = ('id', 'order_status', 'created_at', 'total_price', 'user')
    list_display_links = ('id', 'order_status', 'created_at', 'total_price','user')
    list_per_page = 50
    list_filter = ('order_status', 'payment_method', 'promo', 'cancel', )
    search_fields = ('order_status', 'total_price', 'id')
    search_help_text = 'Введите id заказа, id его статуса или итоговую цену'
    readonly_fields = ('user', 'cart', 'created_at', 'payment_method', 'total_price', 'cancel', 'promo',)
    
    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False
    
    def has_change_permission(self, request, obj=None):
        return True

    def get_readonly_fields(self, request, obj=None):
        if obj:
            return ('user', 'cart', 'created_at', 'payment_method', 'total_price', 'cancel', 'promo',)
        else:
            return ('user', 'cart', 'created_at', 'payment_method', 'total_price', 'cancel', 'promo', 'id',)

    def get_form(self, request, obj=None, **kwargs):
        form = super(OrderAdmin, self).get_form(request, obj, **kwargs)
        if obj:
            form.base_fields['order_status'].widget.attrs['readonly'] = True
        return form


@admin.register(ProductFiles)
class ProductFilesAdmin(admin.ModelAdmin):
    list_display = ('id', 'title',)
    list_display_links = ('id', 'title',)
    list_per_page = 50
    inlines = [FilesInline, ]
    search_fields = ('title',)
    search_help_text = 'Введите название файла'