package com.xinou.volunteer.controller;

import com.jfinal.core.Controller;

/**
 * 组织机构操作
 * Created by shizhida on 16/1/13.
 */
public class OrganizationController extends Controller {

    public void getOrganization(){
        render("/resources/organization.html");
    }

}
