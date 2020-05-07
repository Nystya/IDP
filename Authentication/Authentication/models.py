from rest_framework.authtoken.models import Token
from django.contrib.auth.models import AbstractBaseUser, BaseUserManager, PermissionsMixin
from django.db import models, transaction
from django.utils import timezone

from django.conf import settings
from django.db.models.signals import post_save
from django.dispatch import receiver

import grpc
import Proto.api_pb2 as pb
import Proto.api_pb2_grpc as api

profilesChannel = grpc.insecure_channel('profiles_manager:8002')
profilesApi = api.ProfilesStub(profilesChannel)

class UserManager(BaseUserManager):
    def _create_user(self, email, password, user_type, is_staff, is_superuser, **extra_fields):
        if not email:
            raise ValueError('Users must have an email address')
        now = timezone.now()
        email = self.normalize_email(email)
        user = self.model(
            email=email,
            is_staff=is_staff,
            is_active=True,
            is_verified=False,
            is_superuser=is_superuser,
            last_login=now,
            date_joined=now,
            user_type=user_type,
            **extra_fields
        )
        user.set_password(password)
        user.save(using=self._db)
        return user

    def create_user(self, email, password, user_type, **extra_fields):
        return self._create_user(email, password, user_type, False, False, **extra_fields)

    def create_superuser(self, email, password, user_type, **extra_fields):
        user = self._create_user(email, password, user_type, True, True, **extra_fields)
        user.save(using=self._db)
        return user


class User(AbstractBaseUser, PermissionsMixin):
    email = models.EmailField(max_length=254, unique=True)

    is_staff = models.BooleanField(default=False)
    is_superuser = models.BooleanField(default=False)
    is_active = models.BooleanField(default=True)
    is_verified = models.BooleanField(default=False)

    EMPLOYER = 'employer'
    EMPLOYEE = 'freelancer'
    SUPER = 'SU'
    TYPE = [(EMPLOYER, 'Employer'), (EMPLOYEE, 'Employee'), (SUPER, 'Admin')]
    user_type = models.CharField(max_length=10, choices=TYPE)

    last_login = models.DateTimeField(null=True, blank=True)
    date_joined = models.DateTimeField(auto_now_add=True)

    USERNAME_FIELD = 'email'
    EMAIL_FIELD = 'email'
    REQUIRED_FIELDS = ['user_type']

    objects = UserManager()

    def get_absolute_url(self):
        return "/users/%i/" % (self.pk)

    def __str__(self):
        return self.email

    def get_username(self):
        return self.email

    def isEmployee(self):
        return self.user_type == self.EMPLOYEE

    def isEmployer(self):
        return self.user_type == self.EMPLOYER

#########################################################################################
#########################################################################################
#########################################################################################

@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_auth_token(sender, instance=None, created=False, **kwargs):
    if created:
        if not settings.REST_USE_JWT:
            Token.objects.create(user=instance)

@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_employer_profile(sender, instance=None, created=False, **kwargs):
    if created:
        print("Received signal on employer function")
        user = instance
        if user.user_type == user.EMPLOYER:
            print("Trying to create employer profile")
            try:
                employerProfile = pb.EditEmployerProfileRequest(
                    EUID=pb.ID(ID=user.id.__str__()),
                    phone="",
                    lastName="",
                    firstName="",
                )
                profilesApi.CreateEmployerProfile(employerProfile)
            except Exception as e:
                print("Could not create employer profile: ", e)


@receiver(post_save, sender=settings.AUTH_USER_MODEL)
def create_employee_profile(sender, instance=None, created=False, **kwargs):
    if created:
        print("Received signal on freelancer function")
        user = instance
        if user.user_type == user.EMPLOYEE:
            print("Trying to create freelancer profile")
            try:
                freelancerProfile = pb.EditFreelancerProfileRequest(
                    FUID=pb.ID(ID=user.id.__str__()),
                    phone="",
                    lastName="",
                    firstName="",
                    description="",
                )

                profilesApi.CreateFreelancerProfile(freelancerProfile)
            except Exception as e:
                print("Could not create freelancer profile: ", e)

