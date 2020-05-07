from django.urls import path

from .views import GetAllServices, GetAllSkillCategories, GetSkillsByCategory, CreateJob, GetJob, GetJobs, \
    ApplyForJob, GetApplicants, SelectForJob, GetAcceptedFreelancers, GetAcceptedJobs, GetHistory, GetUserHistory, \
    FinishJob, GetEmployerProfile, GetFreelancerProfile, EditFreelancerProfile, EditEmployerProfile, \
    GetSelfEmployerProfile, GetSelfFreelancerProfile

urlpatterns = [
    path("jobs/", CreateJob, name="create_job"),
    path("jobs/q/", GetJobs, name="get_jobs"),
    path("jobs/apply/<str:id>/", ApplyForJob, name="apply_for_job"),
    path("jobs/applicants/<str:id>/", GetApplicants, name="get_job_applicants"),
    path("jobs/select/<str:jid>/<str:fid>/", SelectForJob, name="select_freelancer_for_job"),
    path("jobs/accepted/<str:id>/", GetAcceptedFreelancers, name="get_accepted_freelancers"),
    path("jobs/accepted/", GetAcceptedJobs, name="get_accepted_jobs_for_freelancer"),
    path("jobs/history/", GetHistory, name="get_user_history"),
    path("jobs/history/<str:id>/", GetUserHistory, name="get_user_history"),
    path("jobs/finish/<str:id>/", FinishJob, name="finish_job"),
    path("jobs/<str:id>/", GetJob, name="get_job"),

    path("services/", GetAllServices, name="get_service_categories"),
    path("skillCategories/<str:serviceCategory>/", GetAllSkillCategories, name="get_skill_categories"),
    path("skills/<str:skillCategory>/", GetSkillsByCategory, name="get_skill"),

    path("profiles/e/", GetSelfEmployerProfile, name="get_self_employer_profile"),
    path("profiles/e/edit/", EditEmployerProfile, name="edit_employer_profile"),
    path("profiles/e/<str:id>/", GetEmployerProfile, name="get_employer_profile"),
    path("profiles/f/", GetSelfFreelancerProfile, name="get_self_freelancer_profile"),
    path("profiles/f/edit/", EditFreelancerProfile, name="edit_freelancer_profiles"),
    path("profiles/f/<str:id>/", GetFreelancerProfile, name="get_freelancer_profile"),
]
