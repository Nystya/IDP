from rest_framework.permissions import BasePermission

class IsEmployee(BasePermission):
    """
    Alows access only to authenticated employees.
    """

    def has_permission(self, request, view):
        return  bool(request.user and request.user.is_authenticated) and bool(request.user.isEmployee())


class IsEmployer(BasePermission):
    """
    Alows access only to authenticated employers.
    """

    def has_permission(self, request, view):
        return bool(request.user and request.user.is_authenticated) and bool(request.user.isEmployer())