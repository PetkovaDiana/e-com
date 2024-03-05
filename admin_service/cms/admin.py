from django.contrib import admin

from .forms import CurrentPromotionsAdmin, BlogAdmin

from .models import Banner, Blog, Vacancies, CurrentPromotions, EmailStatic, \
    PaymentMethod, OrderStatus, MetaTags

@admin.register(EmailStatic)
class EmailStaticAdmin(admin.ModelAdmin):
    list_display = ('courier_email', 'car_image', 'cart_image', 'like_image', 'logo_image',)
    list_display_links = ('courier_email', 'car_image', 'cart_image', 'like_image', 'logo_image',)

    def has_delete_permission(self, request, obj=None):
        return False

    def has_add_permission(self, request, obj=None):
        return False


@admin.register(OrderStatus)
class OrderStatusAdmin(admin.ModelAdmin):
    list_display = ('id', 'name',)
    list_display_links = ('id', 'name',)

    def has_delete_permission(self, request, obj=None):
        return False


@admin.register(PaymentMethod)
class PaymentMethodAdmin(admin.ModelAdmin):
    list_display = ('id', 'title', 'description',)
    list_display_links = ('id', 'title', 'description',)
    list_per_page = 30
    search_fields = ('id', 'title',)
    search_help_text = 'Введите id или название'


admin.site.register(CurrentPromotions, CurrentPromotionsAdmin)
admin.site.register(Blog, BlogAdmin)
admin.site.register(Banner)
admin.site.register(Vacancies)
admin.site.register(MetaTags)